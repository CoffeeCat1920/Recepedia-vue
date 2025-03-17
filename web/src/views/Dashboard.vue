<script>
  export default {
    data() {
      return {
        user: null,
        loggedIn: false,
      }
    },    
    async mounted() {
      try {
        const response = await fetch("/api/data/login", {
          credentials: "include",
          method: "GET",
        });

        if (!response.ok) throw new Error("Unauthorized");

        const data = await response.json();

        this.user = data.name;  
        this.loggedIn = data.loggedIn;

      } catch (error) {
        console.error("Error fetching login info:", error);
      }
    }
  }
</script>

<template>
  <div v-if="loggedIn">
    <h1>Hello, {{ user }}!</h1>
    <h2>Your Recipes</h2>
  </div>
</template>
