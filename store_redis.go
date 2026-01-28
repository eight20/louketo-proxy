/*
Copyright 2015 All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"net/url"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
}

// newRedisStore creates a new redis store
func newRedisStore(location *url.URL) (storage, error) {
	// step: get any password
	password := ""
	if location.User != nil {
		password, _ = location.User.Password()
	}

	// step: parse the url notation
	client := redis.NewClient(&redis.Options{
		Addr:     location.Host,
		DB:       0,
		Password: password,
	})

	return redisStore{
		client: client,
	}, nil
}

// Set adds a token to the store
func (r redisStore) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a token from the store
func (r redisStore) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

// Delete remove the key
func (r redisStore) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

// Close closes of any open resources
func (r redisStore) Close() error {
	if r.client != nil {
		return r.client.Close()
	}

	return nil
}
