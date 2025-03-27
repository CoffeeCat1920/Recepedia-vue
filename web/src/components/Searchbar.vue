<script setup>

  import { useRouter } from 'vue-router';
  import { ref, watch } from 'vue';

  const router = useRouter();

  const props = defineProps({
    searchTerm: {
      type: String,
      default: ''   
    }   
  });

  // Create a local ref to modify the value
  const localSearchTerm = ref(props.searchTerm);

  // Emit changes to the parent
  const emit = defineEmits(['update:searchTerm']);

  watch(localSearchTerm, (newValue) => {
    emit('update:searchTerm', newValue);
  });

  const search = () => {
    router.push({
      path: '/search',
      query: { "searchTerm": localSearchTerm.value }
    });
    console.log('pushing search ')
  };

</script>

<template>

  <div class="flex flex-col items-center text-center">

    <form class="flex flex-col mt-20" @submit.prevent='search' >

      <input class="f-input" name="search" type="text" placeholder="Search" aria-label="Search" v-model="localSearchTerm">

      <div class="relative inline-block mt-10">
        <button type="submit" class="taped-button">Search</button> 
        
        <div class="tape tape-top"></div>
        <div class="tape tape-bottom"></div>
      </div>

    </form>

  </div>

</template>
