package parser

import (
	"bytes"
	"slices"
	"testing"
)

func TestParseLinks(t *testing.T) {
	cases := map[string]struct {
		html string
		want []Link
	}{
		"parsed successfully when single anchor element": {
			html: `<a href="/blog">Blog</a>`,
			want: []Link{
				{
					Text: "Blog",
					Href: "/blog",
				},
			},
		},
		"parse successfully when multiple anchor elements": {
			html: `<a href="/blog">Blog</a><a href="#">Home</a>`,
			want: []Link{
				{
					Text: "Blog",
					Href: "/blog",
				},
				{
					Text: "Home",
					Href: "#",
				},
			},
		},

		"parse successfully when nested elements within anchor tag": {
			html: `<a href="/blog">Blog <span>#1</span></a>`,
			want: []Link{
				{
					Text: "Blog #1",
					Href: "/blog",
				},
			},
		},

		"parse successfully when wrapped in parents": {
			html: `
			<html>

				<body>
					<h1>Hello!</h1>
					<a href="/other-page">A link to another page</a>
				</body>

			</html>`,
			want: []Link{
				{
					Text: "A link to another page",
					Href: "/other-page",
				},
			},
		},

		"parse successfully when comment in text": {
			html: `
			<a href="/blog">Blog <span>page</span><!-- commented text SHOULD NOT be included! --></a>`,
			want: []Link{
				{
					Text: "Blog page",
					Href: "/blog",
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			reader := bytes.NewBufferString(test.html)
			got, err := ParseLinks(reader)
			if err != nil {
				t.Errorf("did not expect error but go one : %v", err)
			}

			if ok := slices.Equal(got, test.want); !ok {
				t.Fatalf("got %v want %v", got, test.want)
			}
		})
	}

}

func TestLinkString(t *testing.T) {
	cases := map[string]struct {
		link Link
		want string
	}{
		"simple link": {
			link: Link{
				Text: "Blog",
				Href: "/blog",
			},
			want: `[Text: "Blog" Href: "/blog"]`,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			got := test.link.String()
			if got != test.want {
				t.Fatalf("got %v want %v", got, test.want)
			}
		})
	}
}
