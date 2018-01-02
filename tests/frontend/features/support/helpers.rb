# Clicking the button fails if it's not in view, so we need to
# scroll to it, as per this answer:
# https://stackoverflow.com/questions/17623075/auto-scroll-
# a-button-into-view-with-capybara-and-selenium
def scroll_into_view(element)
  script = <<-JS
    arguments[0].scrollIntoView(true);
  JS

  page.execute_script(script, element.native)
end

def setup_repo(template=REPO_TEMPLATE_PATH)
  FileUtils.rm_rf(REPO_PATH)
  FileUtils.cp_r(template, REPO_PATH)

  Git.init(REPO_PATH).tap do |g|
    g.add(all: true)
    g.commit("Initial commit")
  end
end

def prevent_modal_animations
  page.execute_script('$(".modal").removeClass("fade")')
end
