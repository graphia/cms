<template>
	<div class="form-group">

		<label for="document">Document Identifier</label>

		<div class="input-group">

			<span class="input-group-addon">
				<label for="custom-document-identifier" class="sr-only">Manually set the document's identifier</label>
				<input name="custom-document-identifier" type="checkbox" v-model="cdi" title="Toggle custom document"/>
			</span>

			<!-- disable tabindex when custom document is disabled -->
			<input	:readonly="!cdi"
					:tabindex="!cdi ? '-1' : '0'"
					name="document"
					class="form-control document form-control-label"
					type="text"
					v-model="documentName"
					required
			/>

			<!-- only display language part if language isn't default, currently hardcoded to 'en' -->
			<span class="input-group-addon extension-indicator">
				index<span v-if="translationEnabled && document.language != 'en'">.{{ document.language }}</span>.md
			</span>

		</div>
	</div>
</template>

<script lang="babel">
	import Accessors from '../../Mixins/accessors';
	import slugify from '../../../javascripts/utilities/slugify.js';

	export default {
		name: "Filename",
		data() {
			return {
				cdi: false, // custom document identifier toggle
				dn: "",     // helper variable used to hold the slugged value
			};
		},
		computed: {
			documentName: {
				get() {
					if (this.cdi) {
						return slugify(this.dn);
					};

					return slugify(this.document.title);
				},
				set(value) {
					let slug = slugify(value);
					this.dn = slug;
				}
			},
			translationEnabled() {
				return this.$store.state.translationEnabled;
			},
		},
		mixins: [Accessors],
		watch: {
			documentName() {
				let slug = slugify(this.documentName);
				this.document.document = slug;
				this.dn = slug;
			}
		}
	};
</script>

<style lang="scss" scoped>
	input.document {
		text-align: right;
	}
</style>