package api

import "github.com/pksep/location_search_server/internal/services/slides"


type SlideHandler struct {
	service *slides.SlidesService
}

func NewSlideHandler(service *slides.SlidesService) *SlideHandler {
	return &SlideHandler{service: service}
}

