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

$application = new Symfony\Component\Console\Application();
$application->setName('CCSwitch');
$application->setVersion('1.0.0');

$application->addCommand(new CCSwitch\Commands\InitCommand());
$application->addCommand(new CCSwitch\Commands\ShowCommand());
$application->addCommand(new CCSwitch\Commands\ListCommand());
$application->addCommand(new CCSwitch\Commands\SetCommand());

$exitCode = $application->run();

exit($exitCode);
