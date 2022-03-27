<script>
	import { session } from '$app/stores';
	import { goto } from '$app/navigation';
	import { post } from '$lib/utils.js';

  let username = '';
	let password = '';
	let errors = null;

  async function doLogin(event) {
		const response = await post(`api/v1/login`, { username, password });

		// TODO handle network errors
		errors = response.errors;

		if (response.user) {
			$session.user = response.user;
			goto('/');
		}
	}
</script>

<svelte:head>
	<title>ApiSim Login</title>
</svelte:head>

<div class="col-md-6 offset-md-2">
    <h1>Login</h1>
    <form class="form-horizontal" on:submit|preventDefault={doLogin}>
        <div class="form-group row">
          <label class="col-sm-2 col-form-label" for="username">Username</label>
          <div class="col-sm-8">
            <input type="text" id="username" name="username" class="form-control" placeholder="Username" bind:value={username} />
            <span class="help-block text-danger"></span>
          </div>
        </div>

        <div class="form-group row">
          <label class="col-sm-2 col-form-label" for="password">Password</label>
          <div class="col-sm-8">
            <input type="password" id="password" name="password" class="form-control" placeholder="Password" bind:value={password} />
            <span class="help-block text-danger"></span>
          </div>
        </div>

        <div class="form-group row">
          <label class="col-sm-2 col-form-label" for="remember">Remember</label>
          <div class="col-sm-8">
            <input type="checkbox" id="remember" name="remember" value="true" />
            <span class="help-block text-danger"></span>
          </div>
        </div>

        <div class="form-group row">
          <div class="col-sm-8">
            <button type="submit" class="btn btn-primary">Login</button>
            <a href="/register" type="submit" class="btn btn-default">Register</a>
          </div>
        </div>
    </form>
</div>
