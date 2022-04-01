<script context="module">
  import * as api from '$lib/api.js';

    export async function load({params, url}) {
		const { id } = params;

        const response = await api.get(`api/v1/apikeys/`+id);

    if (response.success == true) {
        return {
            props: { apikey: response.data }
        }
    } else {
        console.log(response.message)
        return {}
    }

    }
</script>


<script>
    import { page } from '$app/stores'
    const secret = $page.url.searchParams.get('secret')

    export let apikey;
</script>

<svelte:head>
	<title>ApiKey details</title>
</svelte:head>


<div>
    <h1>ApiKey details</h1>

	{#if apikey}
    <div>
        <small>
          Created on <i>{apikey.created_at}</i> |
          <a href="/delete" onclick="return confirm('Are you sure?')">Delete</a>
        </small>
      </div>
      <div>
        Provider: {apikey.provider} <br />
        Name: {apikey.name} <br />
        AccessID: {apikey.access_id} <br />
        {#if secret}
          <span class="text-danger">Please note that this secret is only shown once. Keep it safe</span><br />
          Access Secret: {secret}<br />
        {/if}
        {#if apikey.dlr_url}
          DlrUrl: {apikey.dlr_url}
          <a href="/delete">Delete</a>
          <br />
        {:else}
        <h4>Add dlr-url</h4>
        <form action="/xxx" id="apiKey-add" method="POST" class="form-horizontal" role="form">
          <div class="form-group row">
            <label class="col-sm-2 col-form-label" for="dlr_url">DlrURL</label>
            <div class="col-sm-10">
              <input type="text" id="dlr_url" name="dlr_url" class="form-control" placeholder="Callback Url" value="" />
              <span class="help-block text-danger"></span>
            </div>
          </div>
          <div class="form-group row">
            <div class="col-sm-10">
              <button type="submit" class="btn btn-primary">Add</button>
            </div>
          </div>
        </form>
        {/if}
      </div>
    {:else}
      <p>No Api key details</p>
    {/if}

</div>
