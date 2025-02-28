package slide_test

import (
	"context"
	"testing"

	"github.com/OpenSlides/openslides-autoupdate-service/internal/projector"
	"github.com/OpenSlides/openslides-autoupdate-service/internal/projector/slide"
	"github.com/OpenSlides/openslides-autoupdate-service/pkg/dsmock"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	s := new(projector.SlideStore)
	slide.User(s)

	userSlide := s.Get("user")
	assert.NotNilf(t, userSlide, "Slide with name `user` not found.")

	for _, tt := range []struct {
		name   string
		data   map[string]string
		expect string
	}{
		{
			"Only Username",
			map[string]string{
				"user/1/username":          `"jonny123"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"jonny123"}`,
		},
		{
			"Only Firstname",
			map[string]string{
				"user/1/first_name":        `"Jonny"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"Jonny"}`,
		},
		{
			"Only Lastname",
			map[string]string{
				"user/1/last_name":         `"Bo"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"Bo"}`,
		},
		{
			"Firstname Lastname",
			map[string]string{
				"user/1/first_name":        `"Jonny"`,
				"user/1/last_name":         `"Bo"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"Jonny Bo"}`,
		},
		{
			"Title Firstname Lastname",
			map[string]string{
				"user/1/title":             `"Dr."`,
				"user/1/first_name":        `"Jonny"`,
				"user/1/last_name":         `"Bo"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"Dr. Jonny Bo"}`,
		},
		{
			"Title Firstname Lastname Username",
			map[string]string{
				"user/1/username":          `"jonny123"`,
				"user/1/title":             `"Dr."`,
				"user/1/first_name":        `"Jonny"`,
				"user/1/last_name":         `"Bo"`,
				"user/1/structure_level_$": `["1"]`,
			},
			`{"user":"Dr. Jonny Bo"}`,
		},
		{
			"Title Firstname Lastname Username Level",
			map[string]string{
				"user/1/username":           `"jonny123"`,
				"user/1/title":              `"Dr."`,
				"user/1/first_name":         `"Jonny"`,
				"user/1/last_name":          `"Bo"`,
				"user/1/structure_level_$":  `["1"]`,
				"user/1/structure_level_$1": `"Bern"`,
			},
			`{"user":"Dr. Jonny Bo (Bern)"}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			closed := make(chan struct{})
			defer close(closed)
			ds := dsmock.NewMockDatastore(closed, tt.data)

			p7on := &projector.Projection{
				ContentObjectID: "user/1",
			}

			bs, keys, err := userSlide.Slide(context.Background(), ds, p7on)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expect, string(bs))
			expectedKeys := []string{
				"user/1/username",
				"user/1/title",
				"user/1/first_name",
				"user/1/last_name",
				"user/1/structure_level_$",
				"user/1/structure_level_$1",
			}
			assert.ElementsMatch(t, keys, expectedKeys)
		})
	}
}
