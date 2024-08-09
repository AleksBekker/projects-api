package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encode[T any](w http.ResponseWriter, status int, item T) error {
    data, err := json.Marshal(item)
    if err != nil {
        return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    fmt.Fprint(w, string(data))
    return nil
}

func decode[T any](r *http.Request) (T, error) {
    var item T
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        return item, err
    }
    return item, nil
}
