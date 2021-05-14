package main

//go:generate vugugen -r -skip-go-mod -skip-main ./webapp/components
//go:generate vugugen -r ./webapp/views
//go:generate vgrgen -r ./webapp/views
