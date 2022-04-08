package harbor

import "github.com/lie-inthesun/remotescan/registry"

func (c Client) Image(image string) registry.ImageCli {
	c.image = image
	return c
}
