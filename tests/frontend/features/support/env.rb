require 'capybara'
require 'capybara/cucumber'
require 'capybara/rspec'
require 'cucumber'
require 'fileutils'
require 'git'
require 'net/http'
require 'pry'
require 'pry-byebug'
require 'selenium-webdriver'

REPO_PATH = '../tmp/repositories/cucumber'
REPO_TEMPLATE_PATH = '../backend/repositories/multiple_filetypes'
PID_PATH = '../tmp/cucumber-browser.pid'
DB_PATH = '../../db/cucumber.db'
SAMPLES_PATH = '../backend/samples'
WEB_SERVER_START_ATTEMPTS = 10 # number of 0.1 second delays to wait for server
DOWNLOAD_DIR = "/tmp/downloads" # needs to be in /tmp because permissions 🚨

DRIVER_PREFS = {
  download: {
    prompt_for_download: false,
    directory_upgrade: true,
    default_directory: DOWNLOAD_DIR
  }
}

Capybara.register_driver(:headless_chrome) do |app|

  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome",
        args: %w{--headless --no-sandbox --disable-gpu},
        prefs: DRIVER_PREFS
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
        binary: "/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome",
        args: %w{--no-sandbox --disable-gpu},
        prefs: DRIVER_PREFS
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
  if ENV['DRIVER']
    c.default_driver = ENV['DRIVER'].to_sym
  else
    c.default_driver = :headless_chrome
  end

  #c.default_driver = :firefox
  #c.default_driver = :chrome
  c.app_host = "http://localhost:9095"
end

Before do

  # ensure no repo exists
  FileUtils.rm_rf(REPO_PATH)

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

  # wait for the server to be running before contining
  1.upto(WEB_SERVER_START_ATTEMPTS) do |inc|
    if inc == WEB_SERVER_START_ATTEMPTS
      fail "Timed out waiting for web server to start"
    else
      break if attempt_webserver_call
      sleep 0.1
    end
  end

end

def	attempt_webserver_call
  Net::HTTP.get('localhost', '/', '9095')
  true
rescue Errno::ECONNREFUSED
  false
end

After do
  kill(@pid)
  if FileTest.exist?(DB_PATH)
    File.delete(DB_PATH)
  end
end

def kill(pid)
  Process.kill("HUP", pid)
  Process.wait
  File.delete(PID_PATH)
  @pid = nil
end

World(Capybara)
