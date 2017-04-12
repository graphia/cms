<template>
	<div id="application">

		<nav id="primary">
			<ul>
				<li class="active">primary 1</li>
				<li>primary 2</li>
				<li>primary 3</li>
			</ul>
		</nav>

		<nav id="secondary">
			<ul>
				<li>secondary 1</li>
				<li>secondary 2</li>
				<li>secondary 3</li>
			</ul>
		</nav>

		<div id="router-view-container">
			<div id="router-view">
				<router-view/>
			</div>
		</div>
	</div>
</template>

<script lang="babel">
	console.log("App loaded");
</script>
