<script setup>

  import { useRouter } from 'vue-router';

  const router = useRouter();

  const search = () => {
    router.push("/search");
  };
   
</script>

<template>

  <div class="flex flex-col items-center">

    <form class="flex flex-col mt-20" @submit.prevent='search' >

      <input class="f-input" name="search" type="text" placeholder="Search" aria-label="Search">

      <div class="relative inline-block mt-10">
        <button type="submit" class="taped-button">Search</button> 
        
        <div class="tape tape-top"></div>
        <div class="tape tape-bottom"></div>
      </div>

    </form>

  </div>

</template>
