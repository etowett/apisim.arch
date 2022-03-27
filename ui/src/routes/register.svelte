<script>
    import { session } from '$app/stores';
	import { goto } from '$app/navigation';
	import { post } from '$lib/utils.js';
	import ListErrors from '$lib/ListErrors.svelte';

    let username = '';
    let firstname = '';
    let lastname = '';
	let email = '';
	let password = '';
	let confirmpassword = '';
	let errors = null;

    async function doRegister(event) {

        let request = {
            username, firstname, lastname, email, password, confirmpassword
        }

        console.log(request)

		const response = await post(`api/v1/save`, request);

		// TODO handle network errors
		errors = response.errors;

		if (response.success == true) {
			$session.user = response.data;
			goto('/messages');
		}
	}
</script>
<svelte:head>
    <title>Sign up • Apisim</title>
</svelte:head>

<div class="col-md-6 offset-md-2">
<h1>Register</h1>

<ListErrors {errors}/>

<form class="form-horizontal" on:submit|preventDefault={doRegister}>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="username">Username</label>
    <div class="col-sm-9">
        <input type="text" id="username" name="username" class="form-control" placeholder="Username" bind:value={username} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="firstname">First Name</label>
    <div class="col-sm-9">
        <input type="text" id="firstname" name="firstname" class="form-control" placeholder="First Name" bind:value={firstname} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="lastname">Last Name</label>
    <div class="col-sm-9">
        <input type="text" id="lastname" name="lastname" class="form-control" placeholder="Last Name" bind:value={lastname} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="email">Email</label>
    <div class="col-sm-9">
        <input type="email" id="email" name="email" class="form-control" placeholder="Email" bind:value={email} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="password">Password</label>
    <div class="col-sm-9">
        <input type="password" id="password" name="password" class="form-control" placeholder="Password" bind:value={password} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row ">
    <label class="col-sm-3 col-form-label" for="confirmpassword">Password Confirmation</label>
    <div class="col-sm-9">
        <input type="password" id="confirmpassword" name="confirmpassword" class="form-control" placeholder="Password Confirmation" bind:value={confirmpassword} />
        <span class="help-block text-danger"></span>
    </div>
    </div>

    <div class="form-group row">
    <div class="col-sm-10">
        <button type="submit" class="btn btn-primary">Register</button>
        <a href="/login" type="submit" class="btn btn-default">Login</a>
    </div>
    </div>
</form>
</div>
