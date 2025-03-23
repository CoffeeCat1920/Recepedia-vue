<!-- /src/views/Home.vue -->

<script setup>
  import Brand from '@/components/Brand.vue';
</script>

<template>

  <div class="flex flex-col items-center text-center pt-10 mx-100">

    <div class="card-container">
      
      <div class="card-content px-100">

        <Brand />

        <div class="flex flex-col items-center">

          <form class="flex flex-col mt-20" action="/view/search" >

            <input class="f-input" name="search" type="text" placeholder="Search" aria-label="Search">

            <div class="relative inline-block mt-10">
              <button type="submit" class="taped-button">Search</button> 
              
              <div class="tape tape-top"></div>
              <div class="tape tape-bottom"></div>
            </div>

          </form>

        </div>

      </div>

    </div>

  </div>

  <h2 class="text-center text-amber-900 font-[Monomakh] font-bold text-3xl pt-20">Most Viewed</h2>

</template>
