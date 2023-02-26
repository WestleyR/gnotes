//
//  ViewController.m
//  gnotes-ios
//
//  Created by Westley Rose on 2/23/23.
//

#import "ViewController.h"

#include "ios-lib/gnotes.h"

@interface ViewController ()

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.

    noteTitles = [NSMutableArray new];

    //self.tableViewNotes.dataSource = self;
    //self.tableViewNotes.delegate = self;


    [self.loadingWheelMiddle startAnimating];
    [self.loadingWheelMiddle setHidden:false];

    NSString* tmpDirectory = NSTemporaryDirectory();


    NSString* configFile = [self templateConfigWithTMPDir:tmpDirectory];
    //Download((char*)configFile.UTF8String);

    NSString* jsonIndex = [tmpDirectory stringByAppendingPathComponent:@"notes/index.json"];
    NSString* json = [NSString stringWithContentsOfFile:jsonIndex encoding:NSUTF8StringEncoding error:nil];

    NSLog(@"JSON STRING: %@", json);

    id object = [NSJSONSerialization JSONObjectWithData:[json dataUsingEncoding:NSUTF8StringEncoding] options:0 error:nil];

    id folders = [object valueForKey:@"folders"];

    for (id notes in folders) {
        for (id note in [notes valueForKey:@"notes"]) {
            NSString* title = [note valueForKey:@"title"];
            id isAttachment = [note valueForKey:@"attachment"];
            if ([isAttachment boolValue] == YES) {
                // TODO: attachment titles not working yet...
                title = [notes valueForKey:@"\"attachment_title\""]; // Why double quotes?
            }
            if (title != nil) {
                [noteTitles addObject:title];
            } else {
                NSLog(@"Warning: skipping nil title note: %@", note);
            }
        }
        break; // Only want first index (default Notes)
    }


    //SetupView* setupView = [self.storyboard instantiateViewControllerWithIdentifier:@"setupView"];
    //[self presentModalViewController:setupView animated:YES];



    [self performSegueWithIdentifier:@"SwitchToSetupView" sender:self];

    return;



    //AppDelegate *appDelegate = (AppDelegate *)[[UIApplication sharedApplication] delegate];
    SetupView *aTwoViewController = [self.storyboard instantiateViewControllerWithIdentifier:@"setupView"];

    if (aTwoViewController == nil){
        /* You could use this instead if not using xib:
         yourViewController = [[YourViewController alloc]
         initWithNibName:@"YourViewController"
         bundle:nil];
         */

        UIStoryboard *mainStoryboard = [UIStoryboard storyboardWithName:@"MainStoryboard"
                                                                 bundle: nil];
        aTwoViewController = [mainStoryboard instantiateViewControllerWithIdentifier: @"setupView"];

    }

    // get the view that's currently showing
    UIView *currentView = self.view;
    // get the the underlying UIWindow, or the view containing the current view
    UIView *theWindow = [currentView superview];

    UIView *newView = aTwoViewController.view;

    // remove the current view and replace with myView1
    [currentView removeFromSuperview];
    [theWindow addSubview:newView];

    // set up an animation for the transition between the views
    CATransition *animation = [CATransition animation];
    [animation setDuration:0.5];
    [animation setType:kCATransitionPush];
    [animation setSubtype:kCATransitionFromRight];
    [animation setTimingFunction:[CAMediaTimingFunction functionWithName:kCAMediaTimingFunctionEaseInEaseOut]];

    [[theWindow layer] addAnimation:animation forKey:@"SwitchToView2"];














    /*
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        char* ret = GnotesTest();
        NSLog(@"Got ret from go staic library: %s", ret);

        dispatch_async(dispatch_get_main_queue(), ^{
            [self.loadingWheelMiddle stopAnimating];
            [self.loadingWheelMiddle setHidden:true];
            self.labelResponse.text = [NSString stringWithUTF8String:ret];
        });
    });

     */
}

// Returns a path to a config file
- (NSString*)templateConfigWithTMPDir:(NSString*)tmpDirectory {
    NSString *configPath = [[NSBundle mainBundle] URLForResource:@"config" withExtension:@"ini.tpl"].path;

    NSError* err;
    NSString* template = [NSString stringWithContentsOfFile:configPath encoding:NSUTF8StringEncoding error:&err];
    if (err != nil) {
        NSLog(@"Error reading template file: %@", err);
        return @"";
    }

    NSMutableDictionary* values = [NSMutableDictionary new];

    // Create the cache dir
    [[NSFileManager defaultManager] createDirectoryAtPath:tmpDirectory withIntermediateDirectories:YES attributes:nil error:&err];
    if (err != nil) {
        NSLog(@"Error creating cache dir: %@", err);
        return @"";
    }


    configPath = [tmpDirectory stringByAppendingPathComponent:@"config.ini"];

    // Template values, should not be defined here
    [values setValue:tmpDirectory forKey:@"{{noteDir}}"];
    [values setValue:@"nope" forKey:@"{{accessKey}}"];
    [values setValue:@"nope" forKey:@"{{secretKey}}"];
    [values setValue:@"nope" forKey:@"{{userID}}"];
    [values setValue:@"nope" forKey:@"{{cryptKey}}"];

    for (id key in values) {
        template = [template stringByReplacingOccurrencesOfString:key withString:[values valueForKey:key]];
    }

    NSLog(@"END TEMPLATE:\n%@\n\n", template);
    NSLog(@"END CONFIG PATH: %@", configPath);

    [template writeToFile:configPath atomically:YES encoding:NSUTF8StringEncoding error:&err];
    if (err != nil) {
        NSLog(@"Error writing to file: %@", err);
        return @"";
    }

    return configPath;
}




- (NSInteger)tableView:(nonnull UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return noteTitles.count;
}

- (nonnull UITableViewCell *)tableView:(nonnull UITableView *)tableView cellForRowAtIndexPath:(nonnull NSIndexPath *)indexPath {
    UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"Cell"];
    cell.textLabel.text = noteTitles[indexPath.row];
    return cell;
}

- (void)tableView:(nonnull UITableView*)tableView didSelectRowAtIndexPath:(nonnull NSIndexPath *)indexPath {
    NSLog(@"CLicked row: %@", indexPath);

    // TODO: deselect after downloading note and opening
    [tableView deselectRowAtIndexPath:indexPath animated:YES];
}


@end
