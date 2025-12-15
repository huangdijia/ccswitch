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

use CCSwitch\Profiles;
use Exception;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'list',
    description: 'List all available profiles'
)]
class ListCommand extends Command
{
    protected function configure(): void
    {
        $this
            ->setHelp('This command allows you to list all configured Claude API profiles')
            ->addOption(
                'profiles',
                'p',
                InputOption::VALUE_OPTIONAL,
                'Path to the profiles configuration file',
                getenv('HOME') . '/.ccswitch/ccs.json'
            );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $profilesPath = $input->getOption('profiles');

        try {
            $profiles = new Profiles($profilesPath);
            $defaultProfile = $profiles->default();

            $output->writeln('<info>Available Claude API Profiles:</info>');
            $output->writeln('');

            $profileData = $profiles->data['profiles'] ?? [];
            $descriptions = $profiles->data['descriptions'] ?? [];

            $table = new Table($output);
            $table->setHeaders(['Profile', 'Description', 'URL', 'Model', 'Status']);

            foreach ($profileData as $name => $config) {
                $status = $name === $defaultProfile ? '<info>Default</info>' : '';
                $url = $config['ANTHROPIC_BASE_URL'] ?? '';
                $model = $config['ANTHROPIC_MODEL'] ?? '';
                $description = $descriptions[$name] ?? '';

                $table->addRow([
                    $name,
                    $description,
                    $url,
                    $model,
                    $status,
                ]);
            }

            $table->setStyle('box');
            $table->render();

            $output->writeln('');
            $output->writeln('<comment>Total profiles: ' . count($profileData) . '</comment>');

            return Command::SUCCESS;
        } catch (Exception $e) {
            $output->writeln('<error>Error: ' . $e->getMessage() . '</error>');
            return Command::FAILURE;
        }
    }
}
