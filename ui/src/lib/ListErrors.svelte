<script>
	export let error;
</script>

{#if error}
  <div class="alert alert-danger alert-dismissible fade show" role="alert">
    <strong>Failed!</strong> Please fix errors marked below.
    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    <p>
      {error}
    </p>
  </div>
{/if}
