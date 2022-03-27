<script>
    import {
        Collapse,
        Navbar,
        NavbarToggler,
        NavbarBrand,
        Nav,
        NavItem,
        NavLink
    } from 'sveltestrap/src';

    let isOpen = false;

    function handleUpdate(event) {
        isOpen = event.detail.isOpen;
    }
</script>

<header>
    <Navbar color="dark" dark expand="md" container>
        <NavbarBrand href="/">Apisim</NavbarBrand>
        <NavbarToggler on:click={() => (isOpen = !isOpen)} />
        <Collapse {isOpen} navbar expand="md" on:update={handleUpdate}>
            <Nav navbar>
                <NavItem>
                    <NavLink href="/messages">Messages</NavLink>
                </NavItem>
            </Nav>
        </Collapse>
        <Collapse {isOpen} navbar expand="md" on:update={handleUpdate}>
            <Nav navbar class="ms-auto">
                <NavItem>
                    <NavLink href="/apikeys">ApiKeys</NavLink>
                </NavItem>
                <NavItem>
                    <NavLink href="/login">Login</NavLink>
                </NavItem>
                <NavItem>
                    <NavLink href="/register">Register</NavLink>
                </NavItem>
            </Nav>
        </Collapse>
    </Navbar>
</header>

<main class="container main-content">
    <slot />
</main>

<hr />
<footer class="container text-muted">
    <p class="float-end"><a href="/">Back to top</a></p>
    <p>&copy; 2021 Citizix, Inc. &middot; <a href="/">Privacy</a> &middot; <a href="/">Terms</a></p>
</footer>

<style>
    .main-content {
        padding: 40px;
        min-height: 400px;
    }
</style>
