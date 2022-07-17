package cache

//A redis (maybe) or in memory cache to hold the ids of users who have already authenticated themsleves
//The cache will allow us to skip a trip to mongodb to check if the user exists in the databse on every request
//the cache will mainly be used for users in a cli session

type Cache struct{}

func New() (*Cache, error) {
	return &Cache{}, nil
}

func (cache *Cache) RegisterUser() {

}

func (cache *Cache) FindUser(id string) bool {
	return true
}
