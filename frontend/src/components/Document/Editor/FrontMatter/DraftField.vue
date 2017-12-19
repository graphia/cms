<template>
	<div class="document-draft form-group">
		<label class="form-control-label">
			<input
				type="checkbox"
				v-model="document.draft"
			/>

			Draft
		</label>
	</div>
</template>

<script lang="babel">
	import Accessors from '../../../Mixins/accessors';

	export default {
		name: "PublishedField",
		mixins: [Accessors]
	};
</script>