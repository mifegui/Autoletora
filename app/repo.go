package main

var currentId int = 0

var eventos Eventos

// Give us some seed data
func init() {

	var e Evento
	e.Event = "buy"
	e.Revenue = 49
	var d Data
	d.Key = "store_name"
	d.Value = "Centauro"
	e.Custom_data = append(e.Custom_data, d)
	RepoCriarEvento(e)
}

//func RepoFindTodo(id int) Todo {
//	for _, t := range todos {
//		if t.Id == id {
//			return t
//		}
//	}
//	// return empty Todo if not found
//	return Todo{}
//}
//
////this is bad, I don't think it passes race condtions
func RepoCriarEvento(e Evento) Evento {
	currentId += 1
	e.Id = currentId
	eventos = append(eventos, e)
	return e
}

//
//func RepoDestroyTodo(id int) error {
//	for i, t := range todos {
//		if t.Id == id {
//			todos = append(todos[:i], todos[i+1:]...)
//			return nil
//		}
//	}
//	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
//}
