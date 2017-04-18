<template>
	<div id="application">

		<!-- Primary Navigation Start -->
		<nav class="navbar navbar-toggleable-md navbar-inverse static-top bg-inverse">
			<button class="navbar-toggler navbar-toggler-right hidden-lg-up" type="button" data-toggle="collapse" data-target="#primary" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>

			<a class="navbar-brand" href="#">Graphia CMS</a>

			<div id="primary" class="collapse navbar-collapse">
				<ul class="navbar-nav mr-auto">
					<li class="nav-item active">Files</li>

					<router-link :to="{name: 'document_index', params: {directory: 'documents'}}" class="nav-link">
						Documents
					</router-link>3

					<li><a class="nav-link" href="#">History</a></li>
					<li><a class="nav-link" href="#">Admin</a></li>

				</ul>
			</div>
		</nav>
		<!-- Primary Navigation End -->

		<!-- Router View Container Start -->
		<div class="container-fluid">
			<router-view/>
		</div>
		<!-- Router View Container End -->

	</div>
</template>

<script lang="babel">
export default {
	name: "GraphiaCMS"
}
</script>