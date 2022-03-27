<script>
	export let errors;
</script>

{#if errors}
  <div class="alert alert-danger alert-dismissible fade show">
    <strong>Failed!</strong> Please fix errors marked below!
    <button type="button" class="close" data-dismiss="alert">&times;</button>
	{#each Object.keys(errors) as key}
	  <li>{key} {errors[key]}</li>
	{/each}
  </div>
{/if}
