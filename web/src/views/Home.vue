<!-- /src/views/Home.vue -->

<script setup>
  import Brand from '@/components/Brand.vue';
  import Recommendation from '@/components/RecipeRecommendation.vue';
  import Searchbar from '@/components/Searchbar.vue';
</script>

<template>

  <div class="flex flex-col items-center text-center pt-10 mx-100">

    <div class="card-container">
      
      <div class="card-content px-100">

        <Brand />

        <Searchbar />


      </div>

    </div>

  </div>

  <Recommendation />

</template>
