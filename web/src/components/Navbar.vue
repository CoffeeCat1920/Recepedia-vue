<!-- /component/Navbar.vue -->

<script>
  import { useAuthStore } from "@/stores/authStore";

  export default {
    setup() {
      const auth = useAuthStore();

      const logout = async () => {
        await auth.logout();
      };

      return { auth, logout };
    },
    async mounted() {
      await this.auth.checkLogin(); 
    },
  };
</script>

<template>

  <nav class="navbar">

    <router-link class="brand" to="/">•ᴗ• Recipidia {{auth.loggedIn}} </router-link>

    <div class="nav-links">

      <ul class="nav-list">

        <li class="left-just"><router-link to="/">Home</router-link></li>
        <li class="left-just"><router-link to="/about">About</router-link></li>

        <template v-if="auth.loggedIn">
          <li class="left-just" ><router-link to="/uploadrecipe">Upload Recipe</router-link></li>
          <li class="left-just" ><router-link to="/dashboard">Dashboard</router-link></li>
          <li class="left-just" >
            <a href="#" @click.prevent="logout">Logout</a>
          </li>
        </template>

        <template v-else>
          <li class="left-just"><router-link to="/signup">SignUp</router-link></li>
          <li class="left-just"><router-link to="/login">LogIn</router-link></li>
        </template>


      </ul>

    </div>

  </nav>

</template>
