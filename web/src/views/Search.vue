<script setup>
  
  import { useRoute } from 'vue-router';
  import { computed } from 'vue'

  import Searchbar from '@/components/Searchbar.vue';
  import SearchPage from '@/components/RecipeSearchBoard.vue' 

  const route = useRoute();
  const searchTerm = computed(() => route.query.searchTerm); 

</script>

<template>
  <div>
    <Searchbar class="pb-7" v-model:searchTerm="searchTerm" />
    <h1 class="text-amber-950 font-monomakh text-center pb-6"> Results for {{searchTerm}} </h1>
    <SearchPage v-model:searchTerm="searchTerm" />
  </div>
</template>
