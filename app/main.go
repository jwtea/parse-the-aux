package main

func main() {
	p := NewParser("aux-addon.lua")

	items := p.GetItems()

	so := NewStoreOptions()
	store := NewStore(so)

	for i, wi := range items {
		store.SaveItem(wi)
	}
}
