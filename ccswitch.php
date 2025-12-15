<?php

declare(strict_types=1);
/**
 * This file is part of huangdijia/ccswitch.
 *
 * @link     https://github.com/huangdijia/ccswitch
 * @document https://github.com/huangdijia/ccswitch/blob/main/README.md
 * @contact  Your name <your-mail@gmail.com>
 */
use CCSwitch\ClaudeSettings;
use CCSwitch\Profiles;

require_once __DIR__ . '/vendor/autoload.php';

(function () {
    $home = $_SERVER['HOME'] ?? ($_SERVER['HOMEDRIVE'] . $_SERVER['HOMEPATH']);
    $profilesPath = $options['profiles'] ?? $home . '/.ccswitch/ccs.json';
    $profiles = new Profiles($profilesPath);

    $settingsPath = $options['settings'] ?? $profiles->getSettingsPath() ?? '~/.claude/settings.json';
    $currentSettings = new ClaudeSettings($settingsPath);

    /** @var array<string, null|string|bool> $options */
    $options = getopt('', [
        'called-by::',
        'help::',
        'profiles::',
        'settings::',
        'reset::',
        ...$profiles->getLongOptions(),
    ]);

    // Display Usage
    if (isset($options['help']) || count($options) === 0) {
        $caller = $options['called-by'] ?? null;
        echo $profiles->getUsage($caller) . PHP_EOL;
        exit(0);
    }

    if (isset($options['reset'])) {
        $currentSettings->env = new stdClass();
        if (isset($currentSettings->model)) {
            unset($currentSettings->model);
        }
        $currentSettings->write();
        echo 'Settings reset to default.' . PHP_EOL;
        exit(0);
    }

    $selected = $profiles->default();

    foreach ($options as $name => $value) {
        if ($profiles->has($name)) {
            $selected = $name;
            break;
        }
    }

    if (! $profiles->has($selected)) {
        throw new Exception('Profile not found: ' . $selected);
    }

    // Print
    echo 'Selected profile: ' . $selected . PHP_EOL;

    $env = $profiles->get($selected);
    $currentSettings->env = $env ?: new stdClass();

    if (! isset($env['ANTHROPIC_MODEL'])) {
        unset($currentSettings->model);
    } else {
        $currentSettings->model = $env['ANTHROPIC_MODEL'] ?? null;
    }

    $currentSettings->write();
})();
