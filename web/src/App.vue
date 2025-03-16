<script>

  import { ref } from 'vue'
  import axios from 'axios'

  export default {

    data() {
      return {
        message: null,
        count: 0,
      }
    },

    mounted() {
      const apiUrl = 'http://localhost:8080/api/message';

      axios.get(apiUrl)
        .then( (response) => {
          this.message = response.data.message;
        })
        .catch( (error) => {
          console.error('Error Fetching data from endpoint cause, ', error);
        });
    }

  };

</script>

<template>

  <div class="header">

    <h1 class="logo">•ᴗ• Recepidia</h1>
    <h2>What to have today?</h2>

    <form class="search-bar" action="/view/search" >
      <input class="search-input" name="search" type="text" placeholder="Search" aria-label="Search">
      <button class="search-button">Search</button>
    </form>

  </div>

  <h2>Most Viewed</h2>

</template>
