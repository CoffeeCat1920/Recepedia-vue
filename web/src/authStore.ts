// src/stores/authStore.js
import { defineStore } from "pinia";
import { useRouter } from "vue-router";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null,
    loggedIn: false,
  }),
  actions: {
    async checkLogin() {
      try {
        const response = await fetch("/api/data/login", { credentials: "include" });
        if (response.ok) {
          const data = await response.json();
          this.user = data.User;
          this.loggedIn = true;
        } else {
          this.user = null;
          this.loggedIn = false;
        }
      } catch (error) {
        console.error("Login check failed", error);
      }
    },
    async login(name : string, password : string) {
      try {

        const router = useRouter();

        console.log("Sending:", JSON.stringify({ name: name, password: password }));
        
        const response = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ name, password }),
        });

        if (response.ok) {
          const data = await response.json();
          this.user = data.User;
          this.loggedIn = true;
          router.push('/dashboard');
        } else {
          router.push('/login');
        }
      } catch (error) {
        console.error("Login error", error);
      }

    },
    async logout() {
      const router = useRouter();
      try {
        const response = await fetch("/api/logout", { method: "POST", credentials: "include" });
        if (response.ok) {
          this.user = null;
          this.loggedIn = false;
          router.push('/login')
        } else {
          throw new Error("Logout failed");
        }
      } catch (error) {
        console.error("Logout error", error);
      }
    },
  },
});
