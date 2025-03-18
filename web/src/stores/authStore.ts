// src/stores/authStore.js
import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null,
    loggedIn: false,
  }),
  actions: {
    async checkLogin() {
      try {
        const response = await fetch("/api/data/login", { credentials: "include" });
        if (!response.ok) {
          this.user = null;
          this.loggedIn = false;
          return
        }           
        const data = await response.json();
        this.user = data.User;
        this.loggedIn = true;
      } catch (error) {
        console.error("Login check failed", error);
      }
    },
    async login(name : string, password : string) {
      try {

        const response = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ name: name, password: password }),
        });

        if (!response.ok) {
          return
        }

        const data = await response.json();
        this.user = data.name;
        this.loggedIn = true;
      } catch (error) {
        console.error("Login error", error);
      }

    },
    async logout() {
      try {
        const response = await fetch("/api/logout", { method: "POST", credentials: "include" });
        if (response.ok) {
          this.user = null;
          this.loggedIn = false;
        } else {
          throw new Error("Logout failed");
        }
      } catch (error) {
        console.error("Logout error", error);
      }
    },
  },
});
