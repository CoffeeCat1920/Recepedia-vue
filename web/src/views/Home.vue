<!-- /src/views/Home.vue -->

<script setup>
  console.log("Home.vue");
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
