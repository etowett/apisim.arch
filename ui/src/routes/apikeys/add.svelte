<script>
  import { goto } from '$app/navigation';
  import * as api from '$lib/api.js';
	import ListErrors from '$lib/ListErrors.svelte';

  let provider = '';
	let apiname = '';
  let username = '';
  let dlrurl = '';
  let user_id = 1;

  let submitting = false;

	let error = null;

  async function createApiKey() {

    submitting = true;

    let request = {
      provider, apiname, username, dlrurl, user_id
    }

		const response = await api.post(`api/v1/apikeys`, request, false);

    if (response.success == false) {
      error = response.message
    } else {
      goto('/apikeys/' + response.data.apikey.id + '?secret=' + response.data.secret)
    }

  }
</script>

<svelte:head>
	<title>Add ApiKeys</title>
</svelte:head>

<div class="py-5 px-3">
  <h1>Add Api Key</h1>

  <ListErrors {error}/>

  <form class="form-horizontal" on:submit|preventDefault={createApiKey}>
    <div class="form-group row">
      <label class="col-sm-2 col-form-label" for="provider">Provider</label>
      <div class="col-sm-10">
        <select class="custom-select mr-sm-2 form-control" id="provider" name="provider" bind:value={provider}>
          <option selected value="">-</option>
          <option value="at">Africas Talking</option>
          <option value="rm">Route Mobile</option>
        </select>
        <span class="help-block text-danger"></span>
      </div>
    </div>

    <div class="form-group row">
      <label class="col-sm-2 col-form-label" for="apiname">Name</label>
      <div class="col-sm-10">
        <input type="text" id="apiname" name="apiname" class="form-control" placeholder="Name" bind:value={apiname} />
        <span class="help-block text-danger"></span>
      </div>
    </div>

    <div class="form-group row">
      <label class="col-sm-2 col-form-label" for="username">Username</label>
      <div class="col-sm-10">
        <input type="text" id="username" name="username" class="form-control" placeholder="Username" bind:value={username}/>
        <span class="help-block text-danger"></span>
      </div>
    </div>

    <div class="form-group row">
      <label class="col-sm-2 col-form-label" for="dlrurl">DlrURL</label>
      <div class="col-sm-10">
        <input type="text" id="dlrurl" name="dlrurl" class="form-control" placeholder="DlrURL(optional)" bind:value={dlrurl}/>
        <span class="help-block text-danger"></span>
      </div>
    </div>

    <div class="form-group row">
      <div class="col-sm-10">
        <button type="submit" class="btn btn-primary">Add</button>
        <a href="/apikeys" type="submit" class="btn btn-danger">Cancel</a>
      </div>
    </div>
  </form>
</div>
