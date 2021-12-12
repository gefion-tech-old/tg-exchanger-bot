package redisstore

type RedisStore struct {
	uActions UserActionsI
}

type RedisStoreI interface {
	UserActions() UserActionsI
}

func InitRedisStore(uc UserActionsI) RedisStoreI {
	return &RedisStore{
		uActions: uc,
	}
}

func (s *RedisStore) UserActions() UserActionsI {
	if s.uActions != nil {
		return s.uActions
	}

	s.uActions = &UserActions{}

	return s.uActions
}
