<template>
	<div class="col-md-4">

		<div class="card border-primary m-1">

			<div class="card-header">New document</div>
			<div class="card-body">

				<p class="text-muted">
					Creating a new document will add a new file to this directory. By default, it will be in <em>draft</em> mode,
					and will not be published. Once you are happy with the contents and have removed the <em>draft</em> flag the document
					will be published in the standard manner. All navigation links will be updated automatically.
				</p>

				<router-link :to="{name: 'document_new', params: {directory: this.$route.params.directory}}" class="btn m1 btn-primary">
					Create
				</router-link>

			</div>


		</div>

	</div>
</template>

<script lang="babel">

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "NewDocumentButton",
		mixins: [Accessors]
	};

</script>