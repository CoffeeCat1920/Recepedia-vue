// src/stores/authStore.js
import { defineStore } from "pinia";
import type { Router } from "vue-router";

export const useAuthStore = defineStore("auth", {

  state: () => ({
    user: "",
    loggedIn: false,
    adminLoggedIn: false,
  }),

  actions: {
    async checkLogin() {
      try {
        const response = await fetch("/api/data/login", { credentials: "include" });
        if (!response.ok) {
          this.user = "";
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

    async checkAdminLogin() {
      try {
        const response = await fetch("/api/data/admin/login", { credentials: "include" });
        if (!response.ok) {
          this.user = "";
          this.loggedIn = false;
          return
        }

        const data = await response.json();
        this.adminLoggedIn = true;
      } catch (error) {
        console.error("Login check failed", error);
      }
    },

    async login(name: string, password: string, router: Router) {
      try {

        const response = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ name, password }),
        });

        if (!response.ok) {
          router.push('/login');
          return;
        }

        this.user = name;
        this.loggedIn = true;
        router.push('/dashboard');

      } catch (error) {
        console.error("Login error", error);
      }

    },

    async adminLogin(password: string, router: Router) {
      try {
        const response = await fetch("/api/admin/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ password }),
        });

        if (!response.ok) {
          router.push('/adminlogin');
          return;
        }

        this.adminLoggedIn = true;
        router.push('/admindashboard');
      } catch (error) {
        console.error("Login error", error);
      }
    },

    async logout(router: Router) {
      try {
        const response = await fetch("/api/logout", { method: "POST", credentials: "include" });
        if (!response.ok) {
          throw new Error("Logout failed");
        }
        this.user = "";
        this.loggedIn = false;
        router.push('/login')
      } catch (error) {
        console.error("Logout error", error);
      }
    },

  },
});
