package constants

type Source struct {
	Name string
	Icon string
	URL  string
}

func GetDofusPortalsSource() Source {
	return Source{
		Name: "dofus-portals.fr",
		Icon: "https://i.imgur.com/j8p3M2D.png",
		URL:  "https://dofus-portals.fr",
	}
}
