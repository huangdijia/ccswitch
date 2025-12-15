<?php

declare(strict_types=1);
/**
 * This file is part of huangdijia/ccswitch.
 *
 * @link     https://github.com/huangdijia/ccswitch
 * @document https://github.com/huangdijia/ccswitch/blob/main/README.md
 * @contact  Your name <your-mail@gmail.com>
 */
require_once __DIR__ . '/vendor/autoload.php';

$application = new Symfony\Component\Console\Application('CCSwitch', '1.0.0');

// Set default command to show help when no command is provided
$application->setDefaultCommand('help', true);
$application->addCommand(new CCSwitch\Commands\HelpCommand());
$application->addCommand(new CCSwitch\Commands\InitCommand());
$application->addCommand(new CCSwitch\Commands\ShowCommand());
$application->addCommand(new CCSwitch\Commands\ListCommand());
$application->addCommand(new CCSwitch\Commands\UseCommand());
$application->addCommand(new CCSwitch\Commands\ResetCommand());

$exitCode = $application->run();

exit($exitCode);
