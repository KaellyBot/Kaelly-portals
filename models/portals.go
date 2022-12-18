package models

var (
	SourceDofusPortals = Source{
		Name: "dofus-portals.fr",
		Icon: "https://i.imgur.com/j8p3M2D.png",
		Url:  "https://dofus-portals.fr",
	}
)

type Source struct {
	Name string
	Icon string
	Url  string
}
