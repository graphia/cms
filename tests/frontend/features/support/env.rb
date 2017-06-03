require 'capybara'
require 'capybara/cucumber'
require 'capybara/rspec'
require 'cucumber'
require 'fileutils'
require 'git'
require 'open3'
require 'net/http'
require 'pry'
require 'pry-byebug'
require 'selenium-webdriver'

REPO_PATH = '../tmp/repositories/cucumber'
REPO_TEMPLATE_PATH = '../backend/repositories/small'
PID_PATH = '../tmp/cucumber-browser.pid'
DB_PATH = '../../db/cucumber.db'

Capybara.register_driver(:headless_chrome) do |app|
  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome\ Canary.app/Contents/MacOS/Google\ Chrome\ Canary",
        args: %w{--headless --no-sandbox --disable-gpu}
      }
    )
  )
end

Capybara.register_driver(:chrome) do |app|
  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome\ Canary.app/Contents/MacOS/Google\ Chrome\ Canary",
        args: %w{--no-sandbox --disable-gpu}
      }
    )
  )
end

Capybara.register_driver :firefox do |app|
  Capybara::Selenium::Driver.new(
    app,
    browser: :firefox,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.firefox(marionette: false)
  )
end

Capybara.configure do |c|
  c.default_driver = :headless_chrome
  #c.default_driver = :firefox
  #c.default_driver = :chrome
  c.app_host = "http://localhost:9095"
end

Before do
  if FileTest.exist?(DB_PATH)
    File.delete(DB_PATH)
  end

  FileUtils.rm_rf(REPO_PATH)
  FileUtils.cp_r(REPO_TEMPLATE_PATH, REPO_PATH)

  Git.init(REPO_PATH).tap do |g|
    g.add(all: true)
    g.commit("Initial commit")
  end

  # kill existing pid first
  if FileTest.exist?(PID_PATH)
    begin
      @pid = Pathname.new(PID_PATH).read.to_i

      kill(@pid) if Process.getpgid(@pid) && @pid > 0
    rescue Errno::ESRCH
      # already dead 😵
    end
  end

  @pid = Process.spawn(
    [
      "../../graphia-cms",
      "-config ../../config/cucumber.yml",
      "-log-to-file true"
    ].join(" ")
  )

    #command = "../../graphia-cms -config=../../config/cucumber.yml -log-to-file=true "

    #Open3.popen3(command) do |stdin, stdout, stderr, wait_thr|
      # FIXME Negroni's output is still appearing, work out how to suppress it
    #end

  Pathname.new(PID_PATH).write(@pid)
  visit('/')

end

After do
  kill(@pid)
end

def kill(pid)
  Process.kill("HUP", pid)
  Process.wait
  File.delete(PID_PATH)
  File.delete(DB_PATH)
  @pid = nil
end

World(Capybara)