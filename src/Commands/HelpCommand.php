<?php

declare(strict_types=1);
/**
 * This file is part of huangdijia/ccswitch.
 *
 * @link     https://github.com/huangdijia/ccswitch
 * @document https://github.com/huangdijia/ccswitch/blob/main/README.md
 * @contact  Your name <your-mail@gmail.com>
 */

namespace CCSwitch\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'help',
    description: 'Display help information about CCSwitch'
)]
class HelpCommand extends Command
{
    protected function configure(): void
    {
        $this->setHelp('This command displays help information about CCSwitch and its available commands');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln(
            <<<'EOT'
<info>CCSwitch - Claude Configuration Switcher v1.0.0</info>

A command-line tool to manage and switch between different Claude API configurations.

<comment>Available commands:</comment>
  <info>init</info>        Initialize ccswitch configuration
  <info>list</info>        List all available profiles
  <info>show</info>        Show current profile configuration
  <info>use</info>         Switch the active Claude API profile
  <info>reset</info>       Reset Claude settings to default

<comment>For more information about a specific command, use:</comment>
  ccswitch help <command>

<comment>Examples:</comment>
  <info>ccswitch init</info>                    # Initialize with default configuration
  <info>ccswitch init --full</info>            # Initialize with full configuration
  <info>ccswitch list</info>                   # List all profiles
  <info>ccswitch use myprofile</info>          # Switch to a specific profile
  <info>ccswitch use</info>                    # Interactive profile selection
  <info>ccswitch show</info>                   # Show current configuration

<comment>Configuration files:</comment>
  Profiles: <info>~/.ccswitch/ccs.json</info>
  Settings: <info>~/.claude/settings.json</info>

EOT
        );

        return Command::SUCCESS;
    }
}
