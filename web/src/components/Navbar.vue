<!-- /component/Navbar.vue -->

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
          method: 'GET',
        });

        if (!response.ok) throw new Error("Unauthorized");

        const data = await response.json();
        this.user = data.User;
        this.loggedIn = true;

      } catch (error) {}
    }
  }
</script>

<template>

  <nav class="navbar">

    <router-link class="brand" to="/">•ᴗ• Recipidia</router-link>

    <div class="nav-links">

      <ul class="nav-list">

        <li class="left-just"><router-link to="/">Home</router-link></li>
        <li class="left-just"><router-link to="/about">About</router-link></li>

        <template v-if="loggedIn">
          <li class="left-just" ><router-link to="/uploadrecipe">Upload Recipe</router-link></li>
          <li class="left-just" ><router-link to="/dashboard">Dashboard</router-link></li>
          <li class="left-just" ><router-link to="/api/logout">Logout</router-link></li>
        </template>

        <template v-else>
          <li class="left-just"><router-link to="/signup">SignUp</router-link></li>
          <li class="left-just"><router-link to="/login">Login</router-link></li>
        </template>


      </ul>


    </div>

  </nav>

</template>
