// +build wasm

package main

import (
	"fmt"
	"syscall/js"

	"flag"

	"github.com/Xumeiquer/wallets/webapp/components"
	"github.com/Xumeiquer/wallets/webapp/middleware"
	"github.com/Xumeiquer/wallets/webapp/views"
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
)

func main() {

	mountPoint := flag.String("mount-point", "#vugu_mount_point", "The query selector for the mount point for the root component, if it is not a full HTML component")
	flag.Parse()

	fmt.Printf("[APP] INFO :: Entering main(), -mount-point=%q\n", *mountPoint)
	defer fmt.Printf("[APP] INFO :: Exiting main()\n")

	buildEnv, err := vugu.NewBuildEnv()
	if err != nil {
		panic(err)
	}

	renderer, err := domrender.New(*mountPoint)
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	rootBuilder := VuguSetup(buildEnv, renderer.EventEnv())

	for ok := true; ok; ok = renderer.EventWait() {

		buildResults := buildEnv.RunBuild(rootBuilder)

		err = renderer.Render(buildResults)
		if err != nil {
			panic(err)
		}
	}

}

// VuguSetup performs UI setup and wiring.
func VuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) vugu.Builder {
	host := js.Global().Get("location").Get("host").String()
	pageMap := views.MakeRoutes().WithRecursive(true).WithClean(true).Map()

	router := vgrouter.New(eventEnv)
	api := middleware.NewAPI(host)
	state := middleware.NewState()
	loader := middleware.NewLoader(state, api)

	buildEnv.SetWireFunc(func(b vugu.Builder) {
		fmt.Printf("[APP] INFO :: Wiring component: %T\n", b)
		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			fmt.Println("[APP] INFO :: Setting Router")
			c.NavigatorSet(router)
		}

		if c, ok := b.(middleware.APISetter); ok {
			fmt.Println("[APP] INFO :: Setting API")
			c.APISet(api)
		}

		if c, ok := b.(middleware.StateSetter); ok {
			fmt.Println("[APP] INFO :: Setting State")
			c.StateSet(state)
		}

		if c, ok := b.(middleware.LoaderSetter); ok {
			fmt.Println("[APP] INFO :: Setting Loader")
			c.LoaderSet(loader)
		}
	})

	root := &components.Root{}
	buildEnv.WireComponent(root)

	// pages - add automatically from generated routes
	for path, inst := range pageMap {
		instBuilder := inst.(vugu.Builder)
		router.MustAddRouteExact(path, vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
			root.Body = instBuilder
		}))
	}

	if router.BrowserAvail() {
		err := router.ListenForPopState()
		if err != nil {
			panic(err)
		}

		err = router.Pull()
		if err != nil {
			panic(err)
		}
	}

	return root
}
