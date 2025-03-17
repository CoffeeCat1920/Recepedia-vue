<template>

  <h1> Upload a Recipe </h1>

  <div class="form-group">
    
    <form action="/api/add-recipe" method="post">

      <div>
        <p>
          <label for="Name">Name</label>
          <br>
          <input type="text" name="name"> 
        </p>
      </div>
      
      <div>
        <p>
          <label for="Content">Write here</label>
          <br>
          <textarea type="text" name="content" rows="10" cols="10"></textarea>
        </p>
      </div>
       
      <div>
        <p>
          <button type="submit">Submit</button>
        </p>
      </div>
      
    </form>

  </div>

</template>
