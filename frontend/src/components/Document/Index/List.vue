<template>

	<div>

		<div v-if="this.documents.length == 0">
			<div class="alert alert-info">
				There's nothing here yet
			</div>
		</div>

		<div v-else class="row document-list">

			<div class="col-md-4" v-for="(d, base, i) in groupedTranslations" :key="i">

				<div class="card document-entry m-1" :data-filename="base" :class="{'border-warning': primary(d).draft}" :data-draft="primary(d).draft">

					<div class="card-header">

						<router-link :data-filename="primary(d).document" :to="{name: 'document_show', params: {directory: directoryPath, document: primary(d).document}}">
							{{ primary(d).title || primary(d).filename }}
						</router-link>

						<span v-if="primary(d).draft" class="badge badge-sm badge-warning text-right">
							Draft
						</span>

					</div>

					<div class="card-body">
						<p class="card-text">{{ primary(d).synopsis || description_placeholder }}</p>
					</div>

					<div class="card-footer" v-if="translationEnabled && d.length > 1">
						<ul class="list-inline">
							<li class="list-inline-item" v-for="(t, k) in translations(d)" :key="k" :data-lang="t.languageInfo.name">

								<router-link :to="{name: 'document_show', params: {directory: t.path, filename: t.filename, document: t.document, language_code: t.language}}">
									{{ (t.languageInfo && t.languageInfo.flag) || "missing" }}
								</router-link>
							</li>
						</ul>
					</div>

				</div>

			</div>

		</div>

	</div>


</template>]

<script lang="babel">

	import NewButton from '../Buttons/New';

	export default {
		name: "IndexList",
		props: ["documents", "directoryPath"],
		components: {NewButton},
		computed: {
			groupedTranslations() {

				// FIXME finding translations using the number of dots is potentially
				// a bit fragile, should use a regexp to check for filenames ending
				// in ".xx.md"
				return this
					.documents
					.sort((a,b) => {
						// default language files first
						return (a.filename.split(".").length - b.filename.split(".").length)
					})
					.reduce((summary, doc) => {
						// use the file's basename to group translations
						// let base = doc.filename.split(".")[0]
						let base = doc.document;

						summary[base] ? summary[base].push(doc) : summary[base] = [doc];

						return summary;
					}, {});
			},
			translationEnabled() {
				return this.$store.state.server.translationInfo.translationEnabled;
			},
		},
		methods: {
			primary(files) {
				return files[0];
			},
			translations(files) {
				return files
					.filter((file) => { return file.isTranslation() })
			}
		}
	};
</script>