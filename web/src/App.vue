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

  <div class='app'>
    <h1>{{message}}</h1>
    <button v-on:click="count++">Counter: {{count}}</button>
  </div class='app'>

</template>
