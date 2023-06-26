package models

import "github.com/ditrit/badaas/badorm"

// TODO que se genere automaticamente todo el archivo
// en el unico caso que no se puede generar es el de pointer sin fk
// TODO que pasa si el fk no es un id?
// ahi tampoco lo puedo terminar de determinar
// porque no tengo el id.IsNil()
// para strings se podria hacer comparando con el string vacio
// el resto bueno, que lo determine el usuario sino
// el problema es que para que vaya a string vacia tengo que hacer el truquito del null mapearlo a string vacio,
// lo cual no puedo hacer en todos los casos
// a menos que ya te cree una StringReference que haga el truquito pero chupala
// quien mierda en este mundo usa esas referencias por string
func (m Sale) GetProduct() (*Product, error) {
	return badorm.VerifyStructLoaded[Product](&m.Product)
}

func (m Sale) GetSeller() (*Seller, error) {
	return badorm.VerifyPointerLoaded[Seller](m.SellerID, m.Seller)
	// TODO entonces creo que es mejor siempre trabajar con pointers
	// porque sino en el caso del struct si vos no le pones nada, te crea uno vacio y te lo guarda solo
	//
	// aca puedo mirar el fk (id), que es un *badorm.UUID
	// aunque puedo hacer un preload relations sin hacer el preloadAttributes
	// TODO entonces el preload de una relacion me tiene que asegurar el preload del fk
}

func (m Phone) GetBrand() (*Brand, error) {
	return badorm.VerifyStructLoaded[Brand](&m.Brand)
}

func (m Seller) GetCompany() (*Company, error) {
	return badorm.VerifyPointerLoaded[Company](m.CompanyID, m.Company)
}

// TODO como hago esto cuando este objeto no tiene el fk
// por ejemplo el country, por ahora podria hacer lo que hice con el product
// porque es un struct, pero la idea es pasar todo a pointer
// para evitar las creaciones sin querer
// entonces solo tengo una relacion que puede ser nil lo cargues o no, y no tengo el fk para verificar
// el tema aca tambien es como aseguro el 1 del lado del country si no lo pongo por struct
// porque si ok al que tiene el fk le pongo not null, pero a este imposible
// entonces la forma tiene que ser por modelo y exponerse a la creacion de uno falso
// aunque esta relacion 1 -> 1 no es tan realista creo, no se da tanto
// supongamos entonces que pudiera ser Country 1 -> 0..1 City
// entonces el fk de city va a ser not null, unique y el city del country va a ser un pointer
// entonces que mierda hago?
// en este caso, podria intercambiar donde pongo el fk
// siendo cuentry el que tiene un fk nullable y un pointer
// por su lado, el city pasaria a tener un un country, y lo puedo sacar por lo de los struct
// pero si la relacion es Country 0..1 -> 0..1 City
// y bueno ahi parece que estoy en la mierda, no hay forma de saberlo en alguna de las dos
// igual creo que puede ser una limitacion de la libreria, no es un caso tan comun creeria

// TODO igual ojo que si se puede mezclar pointers y no pointers
// por ejemplo en value.go se hace eso
// TODO que pasa si me ponen la entidad como pointer pero el id sin pointer
