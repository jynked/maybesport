package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

func (c *Cache) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
    user, err := getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    authHeader := r.Header.Get("Authorization")
    fullUser, err := fetchUserFromMokky(user.ID, authHeader)
    if err != nil {
        http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(fullUser.Favourites)
}

func (c *Cache) AddToFavourites(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uniqueId := vars["uniqueId"]
    if uniqueId == "" {
        http.Error(w, "uniqueId required", http.StatusBadRequest)
        return
    }
    user, err := getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    authHeader := r.Header.Get("Authorization")
    fullUser, err := fetchUserFromMokky(user.ID, authHeader)
    if err != nil {
        http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
        return
    }
    if fullUser.Favourites == nil {
        fullUser.Favourites = []string{}
    }
    for _, id := range fullUser.Favourites {
        if id == uniqueId {
            w.WriteHeader(http.StatusOK)
            return
        }
    }
    fullUser.Favourites = append(fullUser.Favourites, uniqueId)
    if err := updateUserInMokky(fullUser, authHeader); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Cache) RemoveFromFavourites(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uniqueId := vars["uniqueId"]
    if uniqueId == "" {
        http.Error(w, "uniqueId required", http.StatusBadRequest)
        return
    }
    user, err := getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    authHeader := r.Header.Get("Authorization")
    fullUser, err := fetchUserFromMokky(user.ID, authHeader)
    if err != nil {
        http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
        return
    }
    if fullUser.Favourites == nil {
        fullUser.Favourites = []string{}
    }
    newFavs := []string{}
    for _, id := range fullUser.Favourites {
        if id != uniqueId {
            newFavs = append(newFavs, id)
        }
    }
    fullUser.Favourites = newFavs
    if err := updateUserInMokky(fullUser, authHeader); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Cache) GetFavouriteItems(w http.ResponseWriter, r *http.Request) {
    user, err := getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    authHeader := r.Header.Get("Authorization")
    fullUser, err := fetchUserFromMokky(user.ID, authHeader)
    if err != nil {
        http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
        return
    }
    if len(fullUser.Favourites) == 0 {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode([]ItemFlatten{})
        return
    }
    allItems := c.GetFlatItems()
    favouriteItems := []ItemFlatten{}
    for _, favId := range fullUser.Favourites {
        for _, item := range allItems {
            if item.UniqueId == favId {
                favouriteItems = append(favouriteItems, item)
                break
            }
        }
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(favouriteItems)
}

func getAuthenticatedUser(r *http.Request) (*User, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return nil, fmt.Errorf("no auth header")
    }
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "https://3b7b2b24dfd8c527.mokky.dev/auth_me", nil)
    req.Header.Set("Authorization", authHeader)
    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unauthorized")
    }
    defer resp.Body.Close()
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func fetchUserFromMokky(userID int, authHeader string) (*User, error) {
    url := fmt.Sprintf("https://3b7b2b24dfd8c527.mokky.dev/users/%d", userID)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", authHeader)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("user not found, status %d", resp.StatusCode)
    }
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func updateUserInMokky(user *User, authHeader string) error {
    user.Password = ""
    body, _ := json.Marshal(user)
    url := fmt.Sprintf("https://3b7b2b24dfd8c527.mokky.dev/users/%d", user.ID)
    req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", authHeader)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("update failed with status %d", resp.StatusCode)
    }
    return nil
}