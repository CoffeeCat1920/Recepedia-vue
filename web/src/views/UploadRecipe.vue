<template>

  <div class="flex flex-col text-center items-start justify-center pl-12 h-screen">

    <div class="card-container">

      <div class="card-content">

        <h1 class="title">Upload a Recipe</h1>
        <h2 class="title-secondary">What looks yummy?</h2>

        <form action="/api/add-recipe" method="post" class="flex flex-col gap-4">
          <p class="flex flex-col">
            <label for="name" class="f-label">Name</label>
            <br>
            <input id="name" type="text" name="name" class="f-input" required />
          </p>

          <p class="flex flex-col">
            <label for="content" class="f-label">Write here</label>
            <br>
            <textarea id="content" name="content" rows="10" class="border-2 border-dashed border-amber-700 font-semibold  text-amber-950" required></textarea>
          </p>

          <div class="relative inline-block mt-10">
            <button type="submit" class="taped-button">Submit</button>
            <div class="tape tape-top"></div>
            <div class="tape tape-bottom"></div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
