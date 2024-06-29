package shell

import (
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/tobischo/gokeepasslib"
)

func findEntries(groups []gokeepasslib.Group, path []string, username string) []gokeepasslib.Entry {
	if len(path) == 0 {
		panic("search path cannot be empty")
	}
	currentGroups := groups
	currentPath := path
	for {
		//fmt.Printf("path: %+v\n", currentPath)
		//fmt.Printf("groups: %+v\n", lo.Map(currentGroups, func(group gokeepasslib.Group, idx int) string { return group.Name }))
		if len(currentPath) == 1 {
			return lo.FlatMap(currentGroups, func(group gokeepasslib.Group, idx int) []gokeepasslib.Entry {
				//fmt.Printf("group name: %s", group.Name)
				return lo.Filter(group.Entries, func(entry gokeepasslib.Entry, idx int) bool {
					//fmt.Printf("entry title: %s\n", entry.GetTitle())
					//fmt.Printf("entry username: %s\n", entry.Get("UserName").Value.Content)
					return entry.GetTitle() == currentPath[0] && entry.Get("UserName").Value.Content == username
				})
			})

		}
		groupName := currentPath[0]
		currentGroups = lo.FlatMap(currentGroups, func(group gokeepasslib.Group, idx int) []gokeepasslib.Group {
			return lo.Filter(group.Groups, func(subgroup gokeepasslib.Group, idx int) bool { return subgroup.Name == groupName })
		})
		currentPath = currentPath[1:]
	}
}

func ReadKeepass(path, pwd, id, username string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(pwd)
	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		panic(err)
	}

	db.UnlockProtectedEntries()

	ids := strings.Split(id, "/")[1:] // on élimine le premier élément qui doit être vide si id commence pas /

	entries := findEntries(db.Content.Root.Groups, ids, username)
	if len(entries) > 1 {
		panic("Multiple Keepass entries found")
	}
	if len(entries) == 0 {
		panic("No Keepass entry found")
	}
	return entries[0].GetPassword()
}
