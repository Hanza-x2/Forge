package Forge

type Application interface {
	Create(driver *Driver)
	Render(driver *Driver, delta float32)
	Resize(driver *Driver, width, height float32)
	Destroy(driver *Driver)
}
